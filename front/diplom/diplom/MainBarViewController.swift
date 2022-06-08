import UIKit

struct EventsInfoResponse: Decodable {
    let error_text: String
    let has_error: Bool
    //let message: AllData
}

//
struct MainResponse: Decodable, Encodable {
    let error_text: String
    let has_error: Bool
    let message: [MyEvents]
    
}
struct MyEvents: Decodable, Encodable {
    let title: String
    let address: String
    let tags: String
    let png: String
    let val: Int
}
//

struct mesData: Decodable {
    let message: AllData
}

struct AllData: Decodable {
    let all_events: Int
    let all_tags: Int
    let events: [AllEvents]
}

struct AllEvents: Decodable {
    let title: String
    let address: String
    let tags: String
    let png: String
    let selected: String
}

struct AllTagsResponse: Decodable {
    let error_text: String
    let has_error: Bool
    let message: [AllTags]
}

struct AllTags: Decodable, Encodable {
    let id: Int64
    let title: String
}

struct UserInfoResponse: Decodable {
    let error_text: String
    let has_error: Bool
    let message: UserInfo
}

struct UserInfo: Decodable {
    let login: String
    let email: String
    let name: String
    let age: String
    let f_tags: String
    let uf_tags: String
}

var evetsDataMain = [MyEvents]()
var Myerr: String = ""

class ModelMyData {
    
    var allTags = [AllTags]()
    var selectedTags = [AllTags]()
    
    init() {
        
    }
    
    func loadData(all:[AllTags], my:[AllTags]) {
        setupData(all: all, my: my)
    }
    
    func setupData(all: [AllTags], my: [AllTags]) {
        allTags = all
        selectedTags = my
    }
    
    func deleteSameTags() {
        var i = 0
        for alltag in allTags {
            for untag in selectedTags {
                if alltag.id == untag.id {
                    allTags.remove(at: i)
                    i -= 1
                }
            }
            i += 1
        }
    }
    
    func addElement(index: Int) {
        if allTags.count != 0 {
            selectedTags.append(allTags[index])
            allTags.remove(at: index)
        }
    }
    
    func deleteElement(index: Int) {
        if selectedTags.count != 0 {
            allTags.append(selectedTags[index])
            selectedTags.remove(at: index)
        }
        
    }
}

class ModelDataTags {
    
    var allTags = [AllTags]()
    var unfavoriteTags = [AllTags]()
    
    init() {
        
    }
    
    func loadData(all:[AllTags], my:[AllTags]) {
        setupData(all: all, my: my)
    }
    
    func setupData(all: [AllTags], my: [AllTags]) {
        allTags = all
        unfavoriteTags = my
    }
    
    func deleteSameTags() {
        var i = 0
        for alltag in allTags {
            for untag in unfavoriteTags {
                if alltag.id == untag.id {
                    allTags.remove(at: i)
                    i -= 1
                }
            }
            i += 1
        }
    }
    
    func addElement(index: Int) {
        if allTags.count != 0 {
            unfavoriteTags.append(allTags[index])
            allTags.remove(at: index)
        }
    }
    
    func deleteElement(index: Int) {
        if unfavoriteTags.count != 0 {
            allTags.append(unfavoriteTags[index])
            unfavoriteTags.remove(at: index)
        }
        
    }
}

class MainBarViewController: UITabBarController {
    
    
    @IBOutlet weak var qwerty: UITabBar!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }

}

class MainEvent: UIViewController {
    
    @IBOutlet weak var EventName: UILabel!
    @IBOutlet weak var ImageEvent: UIImageView!
    @IBOutlet weak var AddressEvent: UILabel!
    @IBOutlet weak var TagsEvent: UILabel!
    @IBOutlet weak var ChangeEvent: UIButton!
    @IBOutlet weak var AcceptEvent: UIButton!
    var index: Int = 0
    override func viewDidAppear(_ animated: Bool) {
        if evetsDataMain.count == 0 {
            index = 0
        }
        while evetsDataMain.count == 0 {
            if Myerr != "" {
                break
            }
        }
        Myerr = ""
        if index == evetsDataMain.count - 1{
            ChangeEvent.isEnabled = false
        }
        EventName.text = evetsDataMain[index].title
        print(evetsDataMain[index].title)
        let url = URL(string: evetsDataMain[index].png)
        if let data = try? Data(contentsOf: url!)
        {
            ImageEvent.image = UIImage(data: data)
        }
        AddressEvent.text = evetsDataMain[index].address
        TagsEvent.text = evetsDataMain[index].tags
        print(evetsDataMain[index].val)
    }
    
    override func viewDidDisappear(_ animated: Bool) {
        navigationController?.popViewController(animated: true)
        dismiss(animated: true, completion: nil)
    }
    
    @IBAction func AcceptEvent(_ sender: UIButton) {
        updateFavoriteTags()
        navigationController?.popViewController(animated: true)
        dismiss(animated: true, completion: nil)
    }
    
    @IBAction func ChangeEvent(_ sender: UIButton) {
        if index != evetsDataMain.count - 1 {
            index += 1
        }
        viewDidAppear(false)
    }
    
    func updateFavoriteTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/update_favorite_tags") else { return }
        var request = URLRequest(url: url)
        let encoder = JSONEncoder()
        let param = try! encoder.encode(evetsDataMain[index])
        print(String(data: param, encoding: .utf8)!)

        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        //guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = param
        
        let session = URLSession.shared
        session.dataTask(with: request) { [] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                    }
                    
                } else {
                    print(httpResp.error_text)
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        index = 0
    }
}

class EventsViewController: UIViewController {
    
    @IBOutlet weak var AddButton: UIButton!
    @IBOutlet weak var AllTagsCol: UIPickerView!
    @IBOutlet weak var MyTagsCol: UIPickerView!
    @IBOutlet weak var DelButton: UIButton!
    @IBOutlet weak var ReadyButton: UIButton!
    
    private var ResAllTags = [AllTags]()
    private var ResMyTags = [AllTags]()
    var myData = ModelMyData()
    
    
    override func viewDidAppear(_ animated: Bool) {
        evetsDataMain.removeAll()
        super.viewDidLoad()
        AllTagsCol.tag = 1
        MyTagsCol.tag = 2
        getAllTags()
        ReadyButton.isEnabled = false
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            tabBarController?.selectedIndex = 0
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        } else {
            ReadyButton.isEnabled = true
        }
    }
    
    func updateViewDisplay() {
        myData.deleteSameTags()
        ResAllTags = myData.allTags
        ResMyTags = myData.selectedTags
        print(ResAllTags)
        print(ResMyTags)
        AllTagsCol.dataSource = self
        AllTagsCol.delegate = self
        MyTagsCol.dataSource = self
        MyTagsCol.delegate = self
    }
    
    @IBAction func ReadyButton(_ sender: UIButton) {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/events/event_by_tag") else { return }
        var request = URLRequest(url: url)
        let encoder = JSONEncoder()
        let param = try! encoder.encode(ResMyTags)
        print(String(data: param, encoding: .utf8)!)

        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        request.httpBody = param
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(EventsInfoResponse.self, from: data)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        let mess = try JSONDecoder().decode(MainResponse.self.self, from: data)
                        evetsDataMain = mess.message
                    }
                    
                } else {
                    print(httpResp.error_text)
                    Myerr = "error"
                    DispatchQueue.main.async {
                        navigationController?.popViewController(animated: true)
                        dismiss(animated: true, completion: nil)
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                    
                }
            } catch {
                print(error)
                Myerr = "error"
                DispatchQueue.main.async {
                    navigationController?.popViewController(animated: true)
                    dismiss(animated: true, completion: nil)
                    tabBarController?.selectedIndex = 0
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            }
        }.resume()
    }
    
    @IBAction func DelButton(_ sender: UIButton) {
        let selected = MyTagsCol.selectedRow(inComponent: 0)
        myData.deleteElement(index: selected)
        updateViewDisplay()
    }
    @IBAction func AddButton(_ sender: UIButton) {
        let selected = AllTagsCol.selectedRow(inComponent: 0)
        myData.addElement(index: selected)
        updateViewDisplay()
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        AllTagsCol.tag = 1
        MyTagsCol.tag = 2
        getAllTags()
    }
    
    func getAllTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(AllTagsResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        DispatchQueue.main.async {
                            ResAllTags = httpResp.message
                            myData.loadData(all: ResAllTags, my: ResMyTags)
                            updateViewDisplay()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                }
            } catch {
                print(error)
                DispatchQueue.main.async {
                    tabBarController?.selectedIndex = 0
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            }
        }.resume()
    }
}

extension EventsViewController: UIPickerViewDataSource {
    func pickerView(_ pickerView: UIPickerView, numberOfRowsInComponent component: Int) -> Int {
        if pickerView.tag == 1 {
            return myData.allTags.count
        } else if pickerView.tag == 2 {
            return myData.selectedTags.count
        } else {
            return 0
        }
    }
    
    func numberOfComponents(in pickerView: UIPickerView) -> Int {
        return 1
    }
}

extension EventsViewController: UIPickerViewDelegate {
    func pickerView(_ pickerView: UIPickerView, titleForRow row: Int, forComponent component: Int) -> String? {
        if pickerView.tag == 1 {
            let a = myData.allTags[row]
            return a.title
        } else if pickerView.tag == 2 {
            let a = myData.selectedTags[row]
            return a.title
        } else {
            return ""
        }
    }
}

class ProfileTagsViewController: UIViewController {
    
    @IBOutlet weak var pickerView: UIPickerView!
    @IBOutlet weak var AddTagButton: UIButton!
    @IBOutlet weak var DeleteTagButton: UIButton!
    @IBOutlet weak var UnfavoriteTagsCollection: UIPickerView!
    
    private var ResAllTags = [AllTags]()
    private var ResUnfavTags = [AllTags]()
    var modelData = ModelDataTags()
    
    override func viewDidAppear(_ animated: Bool) {
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        }
        
    }
    
    func updateViewDisplay() {
        modelData.deleteSameTags()
        ResAllTags = modelData.allTags
        ResUnfavTags = modelData.unfavoriteTags
        print(ResAllTags)
        print(ResUnfavTags)
        pickerView.dataSource = self
        pickerView.delegate = self
        UnfavoriteTagsCollection.dataSource = self
        UnfavoriteTagsCollection.delegate = self
    }
    
    @IBAction func DeleteTagButton(_ sender: UIButton) {
        let selected = UnfavoriteTagsCollection.selectedRow(inComponent: 0)
        modelData.deleteElement(index: selected)
        updateViewDisplay()
    }
    
    @IBAction func AddTagButton(_ sender: UIButton) {
        let selected = pickerView.selectedRow(inComponent: 0)
        modelData.addElement(index: selected)
        updateViewDisplay()
    }
    
    func getUnfavoriteTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/get_unfavorite_tags") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(AllTagsResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            ResUnfavTags = httpResp.message
                            modelData.loadData(all: ResAllTags, my: ResUnfavTags)
                            updateViewDisplay()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                }
            } catch {
                print(error)
                DispatchQueue.main.async {
                    tabBarController?.selectedIndex = 0
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            }
        }.resume()
    }
    
    func updateUnfavoriteTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/update_unfavorite_tags") else { return }
        var request = URLRequest(url: url)
        let encoder = JSONEncoder()
        let param = try! encoder.encode(ResUnfavTags)
        print(String(data: param, encoding: .utf8)!)

        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        //guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = param
        
        let session = URLSession.shared
        session.dataTask(with: request) { [] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                    }
                    
                } else {
                    print(httpResp.error_text)
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    func getAllTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(AllTagsResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        DispatchQueue.main.async {
                            ResAllTags = httpResp.message
                            modelData.loadData(all: ResAllTags, my: ResUnfavTags)
                            updateViewDisplay()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                }
            } catch {
                print(error)
                DispatchQueue.main.async {
                    tabBarController?.selectedIndex = 0
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            }
        }.resume()
    }

    
    override func viewDidDisappear(_ animated: Bool) {
        updateUnfavoriteTags()
        navigationController?.popViewController(animated: true)
        dismiss(animated: true, completion: nil)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        pickerView.tag = 1
        UnfavoriteTagsCollection.tag = 2
        getAllTags()
        getUnfavoriteTags()
    }

}

extension ProfileTagsViewController: UIPickerViewDataSource {
    func pickerView(_ pickerView: UIPickerView, numberOfRowsInComponent component: Int) -> Int {
        if pickerView.tag == 1 {
            return modelData.allTags.count
        } else if pickerView.tag == 2 {
            return modelData.unfavoriteTags.count
        } else {
            return 0
        }
    }
    
    func numberOfComponents(in pickerView: UIPickerView) -> Int {
        return 1
    }
}

extension ProfileTagsViewController: UIPickerViewDelegate {
    func pickerView(_ pickerView: UIPickerView, titleForRow row: Int, forComponent component: Int) -> String? {
        if pickerView.tag == 1 {
            let a = modelData.allTags[row]
            return a.title
        } else if pickerView.tag == 2 {
            let a = modelData.unfavoriteTags[row]
            return a.title
        } else {
            return ""
        }
    }
}

class ProfileViewController: UIViewController {

    @IBOutlet weak var profileName: UILabel!
    @IBOutlet weak var profileLogin: UILabel!
    @IBOutlet weak var profileEmail: UILabel!
    @IBOutlet weak var profileAge: UILabel!
    @IBOutlet weak var favoriteTags: UITextView!
    @IBOutlet weak var unfavoriteTags: UITextView!
    @IBOutlet weak var oldPassword: UITextField!
    @IBOutlet weak var newPassword: UITextField!
    @IBOutlet weak var changePasButton: UIButton!
    @IBOutlet weak var errorLabel: UILabel!
    @IBOutlet weak var exitButton: UIButton!
    
    override func viewDidAppear(_ animated: Bool) {
        DispatchQueue.main.async {
            self.newPassword.text = ""
            self.oldPassword.text = ""
            self.errorLabel.text = ""
        }
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            tabBarController?.selectedIndex = 0
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        } else {
            updateUserInfo()
        }
    }
    
    func updateUserInfo() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/me") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(UserInfoResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            profileName.text = httpResp.message.name
                            profileLogin.text = httpResp.message.login
                            profileEmail.text = httpResp.message.email
                            profileAge.text = httpResp.message.age
                            favoriteTags.text = httpResp.message.f_tags
                            unfavoriteTags.text = httpResp.message.uf_tags
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                }
            } catch {
                print(error)
                DispatchQueue.main.async {
                    tabBarController?.selectedIndex = 0
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            }
        }.resume()
        
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }

    @IBAction func changePasButton(_ sender: UIButton) {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/change_pass") else { return }
        var request = URLRequest(url: url)
        let param = ["old_password": oldPassword.text!, "new_password": newPassword.text!]
        print(param)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            self.errorLabel.textColor = UIColor.systemGreen
                            self.errorLabel.text = "пароль успешно изменен"
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        self.errorLabel.textColor = UIColor.systemRed
                        self.errorLabel.text = httpResp.error_text
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    @IBAction func exitButton(_ sender: UIButton) {
        let Token = UserDefaults.standard
        Token.set("", forKey: "Token")
        Token.synchronize()
        tabBarController?.selectedIndex = 0
        
        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
        self.present(newVC, animated: true, completion: nil)
    }

}

class AllEventsTableView: UITableView {

    /*
    // Only override draw() if you perform custom drawing.
    // An empty implementation adversely affects performance during animation.
    override func draw(_ rect: CGRect) {
        // Drawing code
    }
    */

}

class AllEventsTableViewController: UITableViewController {

    @IBOutlet weak var AllEventsLabel: UILabel!
    @IBOutlet weak var AllTagsLabel: UILabel!
    var MyData: [AllEvents] = []
    
    override func viewDidAppear(_ animated: Bool) {
        getAllData()
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        getAllData()

        // Uncomment the following line to preserve selection between presentations
        // self.clearsSelectionOnViewWillAppear = false

        // Uncomment the following line to display an Edit button in the navigation bar for this view controller.
        // self.navigationItem.rightBarButtonItem = self.editButtonItem
    }

    func getAllData() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/events/") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(EventsInfoResponse.self.self, from: data)
                //print(httpResp.message)
                
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        let mess = try JSONDecoder().decode(mesData.self.self, from: data)
                        print(httpResponse)
                        let AllDataEvents = mess.message
                        MyData = AllDataEvents.events
                        DispatchQueue.main.async {
                            AllEventsLabel.text = String(AllDataEvents.all_events)
                            AllTagsLabel.text = String(AllDataEvents.all_tags)
                            self.tableView.reloadData()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                        tabBarController?.viewDidLoad()
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
    }

    
    // MARK: - Table view data source

    override func numberOfSections(in tableView: UITableView) -> Int {
        // #warning Incomplete implementation, return the number of sections
        return MyData.count
    }

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        // #warning Incomplete implementation, return the number of rows
        return 3
    }

    
    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        if (indexPath.row == 0) {
            let cell = tableView.dequeueReusableCell(withIdentifier: "HeaderCell", for: indexPath) as! HeaderTableViewCell
            cell.selectionStyle = UITableViewCell.SelectionStyle.none
            cell.MyData = MyData[indexPath.section]
            return cell
        } else if (indexPath.row == 1) {
            let cell = tableView.dequeueReusableCell(withIdentifier: "PhotoCell", for: indexPath) as! PhotoTableViewCell
            cell.selectionStyle = UITableViewCell.SelectionStyle.none
            cell.MyData = MyData[indexPath.section]
            return cell
        } else {
            let cell = tableView.dequeueReusableCell(withIdentifier: "CommentCell", for: indexPath) as! CommentTableViewCell
            cell.selectionStyle = UITableViewCell.SelectionStyle.none
            cell.MyData = MyData[indexPath.section]
            return cell
        }
    }
    

    /*
    // Override to support conditional editing of the table view.
    override func tableView(_ tableView: UITableView, canEditRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the specified item to be editable.
        return true
    }
    */

    /*
    // Override to support editing the table view.
    override func tableView(_ tableView: UITableView, commit editingStyle: UITableViewCell.EditingStyle, forRowAt indexPath: IndexPath) {
        if editingStyle == .delete {
            // Delete the row from the data source
            tableView.deleteRows(at: [indexPath], with: .fade)
        } else if editingStyle == .insert {
            // Create a new instance of the appropriate class, insert it into the array, and add a new row to the table view
        }
    }
    */

    /*
    // Override to support rearranging the table view.
    override func tableView(_ tableView: UITableView, moveRowAt fromIndexPath: IndexPath, to: IndexPath) {

    }
    */

    /*
    // Override to support conditional rearranging of the table view.
    override func tableView(_ tableView: UITableView, canMoveRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the item to be re-orderable.
        return true
    }
    */

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destination.
        // Pass the selected object to the new view controller.
    }
    */

}

class HeaderTableViewCell: UITableViewCell {

    @IBOutlet weak var Name: UILabel!
    var MyData: AllEvents? {
        didSet {
            if let setEvent = MyData {
                Name.text = setEvent.title
            }
        }
    }
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}

class PhotoTableViewCell: UITableViewCell {
    
    @IBOutlet weak var ImageEvents: UIImageView!
    var MyData: AllEvents? {
        didSet {
            if let setEvent = MyData {
                let url = URL(string: setEvent.png)
                if let data = try? Data(contentsOf: url!)
                {
                    ImageEvents.image = UIImage(data: data)
                }
            }
            
        }
    }
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}

class CommentTableViewCell: UITableViewCell {

    @IBOutlet weak var Tags: UILabel!
    @IBOutlet weak var Address: UILabel!
    @IBOutlet weak var Liked: UILabel!
    var MyData: AllEvents? {
        didSet {
            if let setEvent = MyData {
                Address.text = setEvent.address
                Tags.text = setEvent.tags
                Liked.text = setEvent.selected
            }
        }
    }
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}

